module Data.ServerConfig exposing (ServerConfig, decoder)

import Json.Decode as D exposing (..)
import Json.Decode.Pipeline exposing (..)


type alias ServerConfig =
    { version : String
    , buildDate : String
    , pop3Listener : String
    , webListener : String
    , smtpConfig : SmtpConfig
    }


decoder : Decoder ServerConfig
decoder =
    succeed ServerConfig
        |> required "version" string
        |> required "build-date" string
        |> required "pop3-listener" string
        |> required "web-listener" string
        |> required "smtp-config" smtpConfigDecoder


type alias SmtpConfig =
    { addr : String
    , defaultAccept : Bool
    , acceptDomains : List String
    , rejectDomains : List String
    , defaultStore : Bool
    , storeDomains : List String
    , discardDomains : List String
    }


smtpConfigDecoder : Decoder SmtpConfig
smtpConfigDecoder =
    succeed SmtpConfig
        |> required "addr" string
        |> required "default-accept" bool
        |> optional "accept-domains" (list string) []
        |> optional "reject-domains" (list string) []
        |> required "default-store" bool
        |> optional "store-domains" (list string) []
        |> optional "discard-domains" (list string) []
