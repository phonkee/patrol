/* Signals

Patrol has functionality to create extern plugins to provide some control.
The functionality is described with Pluginer interface. Patrol has also builtin
blank Plugin and it's recommended to subclass this plugin, so if we add extra
functionality to plugin system there will be already written blank implementation
in this plugin so it won't break already implemented plugins.

(
Patrol has implemented plugin system. For this early version plugins have possibility
to intercept message requests, messages so you can write your own plugin for
refusing messages by information from external api (billing).
In future versions patrol plugins will have possibilities to add custom pages.
Patrol has also migrations framework so plugins can make seamless updates to database.
)

e.g.

```go
type ThrottlePlugin struct {
	patrol.Plugin
}
```

plugin signals:
Every custom plugin has possibility to listen to signals and respond to them.

```go
func(tp *ThrottlePlugin) OnEventRequest(event *RawEvent, rw http.ResponseWriter, r *http.Request) error {
	rw.Write([]byte("please calm down and wait."))
	return errors.New("calm down")
}
```

```go
func(tp *ThrottlePlugin) OnEvent(event Event) {
	rw.Write([]byte("please calm down and wait."))
	return errors.New("calm down")
}
```

*/
package signals
