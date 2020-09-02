# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.4.2] 2020-09-02

- Request state while querying channel

## [2.4.1] 2020-08-20

- Change license to BSD-3

## [2.4.0] 2020-07-29

- Added `options` parameter to `MuteUser` & `MuteUsers` methods, to support `Timeout` option for mute expiration

## [2.3.2] 2020-07-20

- Bump lint to the latest

## [2.3.1] 2020-07-20

### Fixed

- Handle offset and limit in query users

## [2.3.0] 2020-06-25

### Added

- ExtraData support to channel

## [2.2.3] 2020-06-06

### Fixed

- Correct comparison in webhook signature validation

## [2.2.2] 2020-04-30

### Fixed

- Bug in how limit/offset were sent when querying channels

### Added

- Ability to see which other users and channels a user has muted

## [2.2.1] 2020-04-20

### Fixed

- Change jwt dependency to properly generate tokens

## [2.2.0] 2020-04-06

### Fixed

- Add missing or correct wrongly named/typed fields in docs

### Added

- Lots of examples in docs
- Added `Version` helper and used to set a header for requests

### Changed

- Started using upstream for easyjson instead of fork to support unknown keys in JSON

## [2.1.0] 2020-01-23

### Added

- Support for hide channels with clear history

## [2.0.2] - 2020-01-22

### Added

- Support for add message when inviting members or adding\removing moderators.

### Changed

- Fixed issue in GET request body

## [2.0.1] - 2019-11-15

### Fixed

- Add version suffix to go module

## [2.0.0] - 2019-11-14

### Changed

- All methods that update a channel, their members and invites now accept a `*Message` parameter to create a system message

## [1.0.0] - 2019-10-31

### Added

- Support for chat channels and types
- Support for messages
- Support for user and device management
- Support for search; user, channel and message
- Support for moderation and push configuration
- Support for send actions
- Support for partial user update
- Support for sending files
- Support for invite members
