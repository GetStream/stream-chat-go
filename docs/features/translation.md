Chat messages can be translated on-demand or automatically, this allows users speaking different languages on the same channel.

### Message Translation Endpoint

This API endpoint translates an existing message to another language. The source language is inferred from the user language or detected automatically by analyzing its text. If possible it is recommended to store the user language, see " **Set user language** " section later in this page.

```go
translated, err := client.TranslateMessage(ctx, message.ID, "fr")
translated.Message.I18n["fr_text"]
```

The endpoint returns the translated message, updates it and sends a **message.updated** event to all users on the channel.

> [!NOTE]
> Only the text field is translated, custom fields and attachments are not included.


### i18n data

When a message is translated, the `i18n` object is added. The `i18n` includes the message text in all languages and the code of the original language.

The i18n object has one field for each language named using this convention `language-code_text`

Here is an example after translating a message from english into French and Italian.

```json
{
  "fr_text": "Bonjour, J'aimerais avoir plus d'informations sur votre produit.",
  "it_text": "Ciao, vorrei avere maggiori informazioni sul tuo prodotto.",
  "language": "en"
}
```

### Automatic translation

Automatic translation translates all messages immediately when they are added to a channel and are delivered to the other users with the translated text directly included.

Automatic translation works really well for 1-1 conversations or group channels with two main languages.

Let's see how this works in practice:

1. A user sends a message and automatic translation is enabled

2. The language set for that user is used as source language (if not the source language will be automatically detected)

3. The message text is translated into the other language in used on the channel by its members

> [!NOTE]
> When using auto translation, it is recommended setting the language for all users and add them as channel members


### Enabling automatic translation

Automatic translation is not enabled by default. You can enable it for your application via API or CLI from your backend. You can also enable auto translation on a channel basis.

```go
// enable auto-translation only for this channel
channel.Update(ctx, map[string]interface{}{"auto_translation_enabled": true}, nil);

// ensure all messages are translated in english for this channel
update := map[string]interface{}{"auto_translation_enabled": true, "auto_translation_language": "en"}
channel.Update(ctx, update, nil);

// auto translate messages for all channels
enabled := true
settings := &AppSettings{AutoTranslationEnabled: &enabled}
client.UpdateAppSettings(ctx, settings)
```

### Set user language

In order for auto translation to work, you must set the user language or specify a destination language for the channel using the `auto_translation_language` field (see previous code example).

```go
// sets the user language
client.UpsertUser(ctx, &User{ID: "userId", Language: "en"})
```

> [!NOTE]
> Messages are automatically translated from the user language that posts the message to the most common language in use by the other channel members.


### Caveats and limits

- Translation is only done for messages with up to 5,000 characters. Blowin' In The Wind from Bob Dylan contains less than 1,000 characters

- Error messages and commands are not translated (ie. /giphy hello)

- When a message is updated, translations are recomputed automatically

- Changing translation settings or user language have no effect on messages that are already translated

- If there are three or more languages being used by channel members, auto-translate will default to the most common language used by the channel members. Therefore, this feature is best suited for groups with a maximum of **two** main languages.

> [!NOTE]
> A workaround to support more than two languages is to use the translateMessage endpoint to store translated messages for multiple languages, and render the appropriate translation depending on the current users language.


### Available Languages

| Language name         | Language code |
| --------------------- | ------------- |
| Afrikaans             | af            |
| Albanian              | sq            |
| Amharic               | am            |
| Arabic                | ar            |
| Azerbaijani           | az            |
| Bengali               | bn            |
| Bosnian               | bs            |
| Bulgarian             | bg            |
| Chinese (Simplified)  | zh            |
| Chinese (Traditional) | zh-TW         |
| Croatian              | hr            |
| Czech                 | cs            |
| Danish                | da            |
| Dari                  | fa-AF         |
| Dutch                 | nl            |
| English               | en            |
| Estonian              | et            |
| Finnish               | fi            |
| French                | fr            |
| French (Canada)       | fr-CA         |
| Georgian              | ka            |
| German                | de            |
| Greek                 | el            |
| Haitian Creole        | ht            |
| Hausa                 | ha            |
| Hebrew                | he            |
| Hindi                 | hi            |
| Hungarian             | hu            |
| Indonesian            | id            |
| Italian               | it            |
| Japanese              | ja            |
| Korean                | ko            |
| Latvian               | lv            |
| Lithuanian            | lt            |
| Malay                 | ms            |
| Norwegian             | no            |
| Persian               | fa            |
| Pashto                | ps            |
| Polish                | pl            |
| Portuguese            | pt            |
| Romanian              | ro            |
| Russian               | ru            |
| Serbian               | sr            |
| Slovak                | sk            |
| Slovenian             | sl            |
| Somali                | so            |
| Spanish               | es            |
| Spanish (Mexico)      | es-MX         |
| Swahili               | sw            |
| Swedish               | sv            |
| Tagalog               | tl            |
| Tamil                 | ta            |
| Thai                  | th            |
| Turkish               | tr            |
| Ukrainian             | uk            |
| Urdu                  | ur            |
| Vietnamese            | vi            |
