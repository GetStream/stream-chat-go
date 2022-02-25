# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [5.1.0](https://github.com/GetStream/stream-chat-go/compare/v5.0.0...v5.1.0) (2022-02-25)


### Features

* add all app settings ([#180](https://github.com/GetStream/stream-chat-go/issues/180)) ([69aea7f](https://github.com/GetStream/stream-chat-go/commit/69aea7f809eb5421d96139dcc12b1d6f51b3e8a5))
* add query with options ([#184](https://github.com/GetStream/stream-chat-go/issues/184)) ([1f56456](https://github.com/GetStream/stream-chat-go/commit/1f56456705dbddf43d589ff6f4419dfa76d4c5f4))
* add quoted message id top options ([#178](https://github.com/GetStream/stream-chat-go/issues/178)) ([6c5a96c](https://github.com/GetStream/stream-chat-go/commit/6c5a96c9ebddc841a297af4936d7c5395a7527be))
* add translate message endpoint ([#179](https://github.com/GetStream/stream-chat-go/issues/179)) ([d463b3a](https://github.com/GetStream/stream-chat-go/commit/d463b3a08e94a8d870a5730733aef766e76dfe1d))
* typed options for query ([#186](https://github.com/GetStream/stream-chat-go/issues/186)) ([bb2463a](https://github.com/GetStream/stream-chat-go/commit/bb2463ad9479dbe077c17c8ef17183491ae0cad9))

## [5.0.0](https://github.com/GetStream/stream-chat-go/compare/v4.0.1...v5.0.0) (2022-02-03)


### ⚠ BREAKING CHANGES

- `PartialUpdateMessage` method has a new signature
- `Truncatechannel` method has a new signature
- `AddMembers` method has a new signature
- `MarkRead` method has a new signature
- `UpdateCommand` method has a new signature
- `PartialUpdateMessage` method has a new signature
- `MuteUser` method has a new signature
- `MuteUsers` method has a new signature
- `FlagUser` method has a new signature
- `ExportUser` method has a new signature
- `DeactivateUser` method has a new signature
- `ReactivateUser` method has a new signature
- `DeleteUser` method has a new signature
- `BanUser` methods have a new signature
- `UnbanUser` methods have a new signature
- `ShadowBan` methods have a new signature

### Features

* add grants to channeltype ([#166](https://github.com/GetStream/stream-chat-go/issues/166)) ([0a1a824](https://github.com/GetStream/stream-chat-go/commit/0a1a8242d61e4a96084d11bebb0f8e92b79b66ea))
* add import endpoint ([#172](https://github.com/GetStream/stream-chat-go/issues/172)) ([1dd3eba](https://github.com/GetStream/stream-chat-go/commit/1dd3eba5beb5a0b6559cda5a90b48bd9e4c0e2db))
* add offset and limit to listimports ([#174](https://github.com/GetStream/stream-chat-go/issues/174)) ([8c5702b](https://github.com/GetStream/stream-chat-go/commit/8c5702b170b7e5322f761a2b5b8f2feed95613ac))
* enhance connection pooling ([#171](https://github.com/GetStream/stream-chat-go/issues/171)) ([a78a42a](https://github.com/GetStream/stream-chat-go/commit/a78a42afed6e928b255364d64a713351d585b50f))
* extend app config with upload configs ([#170](https://github.com/GetStream/stream-chat-go/issues/170)) ([f4466ca](https://github.com/GetStream/stream-chat-go/commit/f4466ca4e506fec9e2e162757000b0a42e438a43))
* full feature parity ([#168](https://github.com/GetStream/stream-chat-go/issues/168)) ([6cac452](https://github.com/GetStream/stream-chat-go/commit/6cac452969917b7cbd2c8cfe34ca149b90046377))
* improved some apis ([#169](https://github.com/GetStream/stream-chat-go/issues/169)) ([ac44302](https://github.com/GetStream/stream-chat-go/commit/ac443025f2e27ddaa7a849e78b12d6e95da56991))
* swappable http client ([#173](https://github.com/GetStream/stream-chat-go/issues/173)) ([328e767](https://github.com/GetStream/stream-chat-go/commit/328e7677ac5c32496c7cedaf7ca51e8fcf2dbed3))

## [4.0.1] 2021-12-23

- Improve conn closing on errors [#164](https://github.com/GetStream/stream-chat-go/pull/164)

## [4.0.0] 2021-12-17

- Add support for hiding history while adding a member [#149](https://github.com/GetStream/stream-chat-go/pull/149)
- Add support for truncate options (hard_delete, truncated_at, system message) [#151](https://github.com/GetStream/stream-chat-go/pull/151)
- Add support for context in every call [#153](https://github.com/GetStream/stream-chat-go/pull/153)
- Add support for exposing API errors [#154](https://github.com/GetStream/stream-chat-go/pull/154)
- Add support for rate limit headers in responses [#156](https://github.com/GetStream/stream-chat-go/pull/156)
- Add support for permissions v2 [#152](https://github.com/GetStream/stream-chat-go/pull/152) [#161](https://github.com/GetStream/stream-chat-go/pull/161)
- Drop import channel messages endpoint support [#155](https://github.com/GetStream/stream-chat-go/pull/155)
- Drop unflag endpoint support [#157](https://github.com/GetStream/stream-chat-go/pull/157)
- Drop update user in favor of upsert user [#158](https://github.com/GetStream/stream-chat-go/pull/158)
- Require go1.16 [#159](https://github.com/GetStream/stream-chat-go/pull/159)

## [3.14.0] 2021-11-17

- Add support for shadow banning user [#148](https://github.com/GetStream/stream-chat-go/pull/148)
  - ShadowBan
  - RemoveShadowBan
- Add support for pinning messages [#148](https://github.com/GetStream/stream-chat-go/pull/148)
  - PinMessage
  - UnPinMessage
- Add support for partial updating messages [#148](https://github.com/GetStream/stream-chat-go/pull/148)
  - PartialUpdateMessage
- Add support for updating channel ownership for Deleted Users [#147](https://github.com/GetStream/stream-chat-go/pull/147)

## [3.13.0] 2021-11-01

- Add support for async endpoints
  - Delete channels
  - Delete users
  - Export channels
- Add support for async url enrichment app configuration
- Remove base url from readme
  - To simplify setup and unnecessary with edge

## [3.12.2] 2021-09-01

- Use edge as base url at default
- Change jwt dependency for security fixes
- Use POST instead of GET in query channels
- Test with go 1.17
  - further details ([#137](https://github.com/GetStream/stream-chat-go/pull/137))

## [3.12.1] 2021-08-19

- Add missing configuration fields to firebase config ([#135](https://github.com/GetStream/stream-chat-go/pull/135))

## [3.12.0] 2021-08-14

- Add support for message hard delete ([#133](https://github.com/GetStream/stream-chat-go/pull/133))

## [3.11.3] 2021-07-30

- Simplify send user custom event payload ([#131](https://github.com/GetStream/stream-chat-go/pull/131))

## [3.11.2] 2021-06-30

- Proxy command list to createTypeRequest when creating channel type ([#130](https://github.com/GetStream/stream-chat-go/pull/130))

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
