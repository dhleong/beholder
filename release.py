#!/usr/bin/env python
#
# Release script for beholder
#

import hashlib
import urllib
from collections import OrderedDict

try:
    from hostage import *  #pylint: disable=unused-wildcard-import,wildcard-import
except ImportError:
    print "!! Release library unavailable."
    print "!! Use `pip install hostage` to fix."
    print "!! You will also need an API token in .github.token,"
    print "!!  a .hubrrc config, or `brew install hub` configured."
    print "!! A $GITHUB_TOKEN env variable will also work."
    exit(1)

#
# Globals
#

notes = File(".last-release-notes")
latestTag = git.Tag.latest()

def sha256(fileUrl, blockSize=65536):
    # based on: https://gist.github.com/rji/b38c7238128edf53a181
    hasher = hashlib.sha256()
    shafp = urllib.urlopen(fileUrl)
    for block in iter(lambda: shafp.read(blockSize), b''):
        hasher.update(block)
    shafp.close()
    return hasher.hexdigest()

def formatIssue(issue):
    return "- {title} (#{number})\n".format(
        number=issue.number,
        title=issue.title)

def buildLabeled(labelsToTitles):
    """Given a set of (label, title) tuples, produces an
    OrderedDict whose keys are `label`, and whose values are
    dictionaries containing 'title' -> `title`, and
    'content' -> string. The iteration order of the dictionary
    will preserve the ordering of the provided tuples
    """
    result = OrderedDict()
    for k, v in labelsToTitles:
        result[k] = {'title': v, 'content': ''}
    return result

def buildDefaultNotes(_):
    if not latestTag: return ''

    logParams = {
        'path': latestTag.name + "..HEAD",
        'grep': ["Fix #", "Fixes #", "Closes #"],
        'pretty': "format:- %s"}
    logParams["invertGrep"] = True
    msgs = git.Log(**logParams).output()

    contents = ''

    lastReleaseDate = latestTag.get_created_date()
    if lastReleaseDate.tzinfo:
        # pygithub doesn't respect tzinfo, so we have to do it ourselves
        lastReleaseDate -= lastReleaseDate.tzinfo.utcoffset(lastReleaseDate)
        lastReleaseDate.replace(tzinfo=None)

    closedIssues = github.find_issues(state='closed', since=lastReleaseDate)

    labeled = buildLabeled([
        ['feature', "New Features"],
        ['enhancement', "Enhancements"],
        ['bug', "Bug Fixes"],
        ['_default', "Other resolved tickets"],
    ])

    if closedIssues:
        for issue in closedIssues:
            found = False
            for label in labeled.iterkeys():
                if label in issue.labels:
                    labeled[label]['content'] += formatIssue(issue)
                    found = True
                    break
            if not found:
                labeled['_default']['content'] += formatIssue(issue)

    for labeledIssueInfo in labeled.itervalues():
        if labeledIssueInfo['content']:
            contents += "\n**{title}**:\n{content}".format(**labeledIssueInfo)

    if msgs: contents += "\n**Notes**:\n" + msgs
    return contents.strip()

#
# Verify
#

verify(Grep("stopship", inDir="src").foundAny(silent=False)) \
    .then(echoAndDie("I don't think so"))

version = verify(File("main.go")
                 .filtersTo(RegexFilter('const beholderVersion = "(.*)"'))
                ).valueElse(echoAndDie("No version!?"))
versionTag = git.Tag(version)

verify(versionTag.exists())\
    .then(echoAndDie("Version `%s` already exists!" % version))

#
# Make sure all the tests pass
#

verify(Execute("go test github.com/dhleong/beholder/src")).succeeds(silent=False).orElse(die())

#
# Build the release notes
#

initialNotes = verify(notes.contents()).valueElse(buildDefaultNotes)
notes.delete()

verify(Edit(notes, withContent=initialNotes).didCreate())\
        .orElse(echoAndDie("Aborted due to empty message"))

releaseNotes = notes.contents()

#
# Compile
#
versions = [
    # (label, os, arch) tuples
    ("macOS", "darwin", "amd64"),
    ("windows-x64", "windows", "amd64"),
]
compiled = []
for (buildLabel, os, arch) in versions:
    f = 'bin/beholder-%s-%s' % (version, buildLabel)
    if os == "windows":
        f += ".exe"

    print "Compiling:", f
    cmd = 'env GOOS=%s GOARCH=%s go build -v -o %s' % (os, arch, f)
    verify(Execute(cmd)).succeeds(silent=False)

    compiled.append(f)

#
# Upload to github
#

print "Uploading to Github..."

verify(versionTag).create()
verify(versionTag).push("origin")

gitRelease = github.Release(version)
verify(gitRelease).create(body=releaseNotes)

for f in compiled:
    print "Uploading", f
    verify(gitRelease).uploadFile(f, 'application/octet-stream')

#
# Update homebrew repo
#

print "Updating homebrew..."

tarUrl = 'https://github.com/dhleong/beholder/archive/%s.tar.gz' % version
tarSha = sha256(tarUrl)

homebrewConfig = github.Config("dhleong/homebrew-beholder")
formulaFile = github.RepoFile("/Formula/beholder.rb", config=homebrewConfig)
oldContents = formulaFile.read()

newContents = oldContents
newContents = re.sub('url "[^"]+"', 'url "%s"' % tarUrl, newContents)
newContents = re.sub('sha256 "[^"]+"', 'sha256 "%s"' % tarSha, newContents)

print "     url <-", tarUrl
print "  sha256 <-", tarSha
commit = 'Update for v%s' % version
verify(formulaFile).write(newContents, commitMessage=commit)

#
# Success! Now, just cleanup and we're done!
#

notes.delete()

print "Done! Published %s" % version
