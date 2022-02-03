const versionFileUpdater = {
    MAJOR_REGEX: /versionMajor = (\d+)/,
    MINOR_REGEX: /versionMinor = (\d+)/,
    PATCH_REGEX: /versionPatch = (\d+)/,

    readVersion: function (contents) {
        const major = this.MAJOR_REGEX.exec(contents)[1];
        const minor = this.MINOR_REGEX.exec(contents)[1];
        const patch = this.PATCH_REGEX.exec(contents)[1];

        return `${major}.${minor}.${patch}`;
    },

    writeVersion: function (contents, version) {
        const splitted = version.split('.');
        const [major, minor, patch] = [splitted[0], splitted[1], splitted[2]];

        return contents
            .replace(this.MAJOR_REGEX, `versionMajor = ${major}`)
            .replace(this.MINOR_REGEX, `versionMinor = ${minor}`)
            .replace(this.PATCH_REGEX, `versionPatch = ${patch}`);
    }
}

const moduleVersionUpdater = {
    GO_MOD_REGEX: /stream-chat-go\/v(\d+)/g,

    readVersion: function (contents) {
        return this.GO_MOD_REGEX.exec(contents)[1];
    },

    writeVersion: function (contents, version) {
        const major = version.split('.')[0];

        return contents.replace(this.GO_MOD_REGEX, `stream-chat-go/v${major}`);
    }
}

module.exports = {
    bumpFiles: [
        { filename: './version.go', updater: versionFileUpdater },
        { filename: './go.mod', updater: moduleVersionUpdater },
        { filename: './README.md', updater: moduleVersionUpdater },
    ],
}
