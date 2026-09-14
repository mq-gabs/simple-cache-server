# Protocol Definition - SCAS

## Summary

The protocol is divided in two main parts: header and payload. The header is composed by a sequence of bytes with fixed length. The header contains important definitions for the command execution.

### Header Structure

| Name           | Example value | Size    | Description                                  |
| -------------- | ------------- | ------- | -------------------------------------------- |
| Magic          | 0x5ca5        | 2 bytes | Fixed initial bytes. Define the protocol     |
| Version        | 0x01          | 1 byte  | Define the used protocol structure/usability |
| Command        | 0x02          | 1 byte  | Instruction to execute                       |
| Flags          | 0x1f          | 1 byte  | Set command behaviors and modifiers          |
| Payload Length | 0x000010ff    | 4 bytes | Define how long is the payload               |

### Payload Structure

It can be any size limited by the payload length section. How the payload is structured depends on the command.
