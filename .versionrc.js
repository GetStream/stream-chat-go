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
            .replace(this.MAJOR_REGEX.exec(contents)[0], `versionMajor = ${major}`)
            .replace(this.MINOR_REGEX.exec(contents)[0], `versionMinor = ${minor}`)
            .replace(this.PATCH_REGEX.exec(contents)[0], `versionPatch = ${patch}`);
    }
}

module.exports = {
    bumpFiles: [{ filename: './version.go', updater: versionFileUpdater }],
}
