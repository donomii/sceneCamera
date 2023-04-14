# Message module

Pass messages between modules

The message module does not use any external libraries, and can be used in any project.

## Overview

The message registry allows communication between modules.  It supports dynamic and changing message handles, and is a solution to the common problems with import cycles and dependencies that other systems have.

The use is simple.  You register a handler for a message, and then send the message.  The handler will be called with the message name, an id, and the arguments.  You provide the id when you register.  The id is used to identify the handler, so you can remove it later.  Registering another handler with the same id will replace the old handler.

The message registry is thread-safe, and can be used from multiple goroutines.  However the handler will run in the same routine as the caller, so be careful about calling graphics functions from the handler.  Many graphics libraries, such as OpenGL, are not thread-safe.

The message system is designed for notifications of events, less for passing large amounts of data, or for large numbers of events, or for complex execution flow.  It works well for simple things like "the user has clicked on this button", or "the user has pressed this key", or "the user has moved the mouse".

Messages are never queued, they are delivered immediately.

## API

### Register( name string, id string, handler func(name , id string, args interface{}))

Register your message handler with the message system.  Name is the lookup name for the message, id is a free text field that you can use to identify your handler, and handler is the function that will be called when the message is sent.

There can be as many handlers as you want for a message.  The id is used to identify the handler, so you can remove it later.

Registering another handler with the same id will replace the old handler.


### func SendMessage( name string, args interface{})

Send a message.  The name is the lookup name for the message, and args is the data to be sent to the handler.  Sendmessage calls the handler immediately (i.e. synchronously)

SendMessage doesn't check types, so if you send the wrong type of data to the handler, it will probably panic.

### func Unregister(name string, id string)

Delete a handler.  The name and id must be the same as the ones used to register the handler.