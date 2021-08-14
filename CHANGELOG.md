# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.12.0] 2021-08-14

- Add support for message hard delete ([#133](https://github.com/GetStream/stream-chat-go/pull/133)

## [3.11.3] 2021-07-30

- Simplify send user custom event payload ([#131](https://github.com/GetStream/stream-chat-go/pull/131)

## [3.11.2] 2021-06-30

- Proxy command list to createTypeRequest when creating channel type ([#130](https://github.com/GetStream/stream-chat-go/pull/130)

## [3.11.1] 2021-06-29

- Update version header passed to server

## [3.11.0] 2021-06-29

- Add missing channel config support ([#129](https://github.com/GetStream/stream-chat-go/pull/129))

## [3.10.0] 2021-06-25

- Support search improvements of backend ([#128](https://github.com/GetStream/stream-chat-go/pull/128))

## [3.9.1] 2021-06-04

- Add missing `channel.created` event ([#127](https://github.com/GetStream/stream-chat-go/pull/127))

## [3.9.0] 2021-05-31

- Add support for query message flags ([#125](https://github.com/GetStream/stream-chat-go/pull/125))
- Add support for app and user level token revoke ([#121](https://github.com/GetStream/stream-chat-go/pull/121), [#126](https://github.com/GetStream/stream-chat-go/pull/126))

## [3.8.0] 2021-05-18

- Add disabled field to channels ([#124](https://github.com/GetStream/stream-chat-go/pull/124))

## [3.7.0] 2021-05-17

- Add user custom events ([#119](https://github.com/GetStream/stream-chat-go/pull/119))
- Use proxy as default base url ([#122](https://github.com/GetStream/stream-chat-go/pull/122))
- Run tests sequentially to prevent conflicting app state ([#120](https://github.com/GetStream/stream-chat-go/pull/120))
- Automatically clean old data from previous failing tests ([#123](https://github.com/GetStream/stream-chat-go/pull/123))

## [3.6.1] 2021-03-10

- Update internally how endpoints are handled for GetRateLimits endpoint ([#117](https://github.com/GetStream/stream-chat-go/pull/117))

## [3.6.0] 2021-03-09

- Fix update channel type endpoint ([#116](https://github.com/GetStream/stream-chat-go/pull/116))
- Add push notifications enable / disable flag for channel types ([#116](https://github.com/GetStream/stream-chat-go/pull/116))

## [3.5.0] 2021-03-08

- Add get rate limit endpoint support ([#115](https://github.com/GetStream/stream-chat-go/pull/115))
- Add replace go 1.14 with go 1.16 in CI

## [3.4.0] 2021-02-22

- Add options to send message to configure its behavior ([#114](https://github.com/GetStream/stream-chat-go/pull/114))

## [3.3.1] 2021-02-09

- Ensure un/mute a channel works without query the channel first ([#113](https://github.com/GetStream/stream-chat-go/pull/113))

## [3.3.0] 2021-01-22

- Add `UpsertUser` and `UpsertUsers`, and deprecate `UpdateUser` and `UpdateUsers` ([#111](https://github.com/GetStream/stream-chat-go/pull/111))
- Bump lint tool and improve godoc

## [3.2.0] 2021-01-18

- Add team into user and channel for multi-tenant ([#110](https://github.com/GetStream/stream-chat-go/pull/110))

## [3.1.0] 2020-12-17

- Add channel partial update ([#109](https://github.com/GetStream/stream-chat-go/pull/109))

## [3.0.3] 2020-12-14

- Fix duration type in channel mute expiration from seconds to milliseconds ([#108](https://github.com/GetStream/stream-chat-go/pull/108))

## [3.0.2] 2020-12-10

- Support zero as message/member limit in query channels

## [3.0.1] 2020-11-10

- Handle member/message limit in query channels ([#106](https://github.com/GetStream/stream-chat-go/pull/106))

## [3.0.0] 2020-09-24

- Drop client/channel interfaces ([#98](https://github.com/GetStream/stream-chat-go/pull/98))
- Receive string in client initialization ([#99](https://github.com/GetStream/stream-chat-go/pull/99))
- Generate string token instead byte slice ([#100](https://github.com/GetStream/stream-chat-go/pull/100))
- Require go1.14 and above ([#101](https://github.com/GetStream/stream-chat-go/pull/101))

## [2.8.0] 2020-09-24

- Add bulk message import into a channel

## [2.7.0] 2020-09-24

- Add custom command endpoints
- Add missing methods of channel interface

## [2.6.1] 2020-09-23

- Handle members better for reserved fields in query members of a channel

## [2.6.0] 2020-09-18

- Add support for query members of a channel

## [2.5.0] 2020-09-18

- Add support for silent messages
- Test go 1.14 and 1.15 in CI

## [2.4.3] 2020-09-17

- Drop easyjson in favor of standard library (not noticeable from client perspective)
- Bump golangci-lint and replace impi with native linter gci

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
