---
title: Logging
linkTitle: Logging
weight: 20
---

# Overview

Operator SDK-generated operators use the [`logr`][godoc_logr] interface to log. This log interface has several backends such as [`zap`][repo_zapr], which the SDK uses in generated code by default. [`logr.Logger`][godoc_logr_logger] exposes structured logging methods that help create machine-readable logs and adding a wealth of information to log records.

## Default zap logger

Operator SDK uses a `zap`-based `logr` backend when scaffolding new projects. To assist with configuring and using this logger, the SDK includes several helper functions.

In the simple example below, we add the zap flagset to the operator's command line flags with `BindFlags()`, and then set the controller-runtime logger with `zap.Options{}`.

By default, `zap.Options{}` will return a logger that is ready for production use. It uses a JSON encoder, logs starting at the `info` level. To customize the default behavior, users can use the zap flagset and specify flags on the command line. The zap flagset includes the following flags that can be used to configure the logger:

* `--zap-devel`: Development Mode defaults(encoder=consoleEncoder,logLevel=Debug,stackTraceLevel=Warn)
			  Production Mode defaults(encoder=jsonEncoder,logLevel=Info,stackTraceLevel=Error)
* `--zap-encoder`: Zap log encoding ('json' or 'console')
* `--zap-log-level`: Zap Level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error',
			       or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")
* `--zap-stacktrace-level`: Zap Level at and above which stacktraces are captured (one of 'info' or 'error')

Consult the controller-runtime [godocs][logging_godocs] for more detailed flag information.

### A simple example

Operators set the logger for all operator logging in `main.go`. To illustrate how this works, try out this simple example:

```Go
package main

import (
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var globalLog = logf.Log.WithName("global")
func main() {
	// Add the zap logger flag set to the CLI. The flag set must
	// be added before calling flag.Parse().
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	logf.SetLogger(logger)

	scopedLog := logf.Log.WithName("scoped")

	globalLog.Info("Printing at INFO level")
	globalLog.V(1).Info("Printing at DEBUG level")
	scopedLog.Info("Printing at INFO level")
	scopedLog.V(1).Info("Printing at DEBUG level")
}
```

#### Output using the defaults
```console
$ go run main.go
INFO[0000] Running the operator locally in namespace default.
{"level":"info","ts":1587741740.407766,"logger":"global","msg":"Printing at INFO level"}
{"level":"info","ts":1587741740.407855,"logger":"scoped","msg":"Printing at INFO level"}
```

#### Output overriding the log level to 1 (debug)
```console
$ go run main.go --zap-log-level=debug
INFO[0000] Running the operator locally in namespace default.
{"level":"info","ts":1587741837.602911,"logger":"global","msg":"Printing at INFO level"}
{"level":"debug","ts":1587741837.602964,"logger":"global","msg":"Printing at DEBUG level"}
{"level":"info","ts":1587741837.6029708,"logger":"scoped","msg":"Printing at INFO level"}
{"level":"debug","ts":1587741837.602973,"logger":"scoped","msg":"Printing at DEBUG level"}
```
## Custom zap logger

In order to use a custom zap logger, [`zap`][controller_runtime_zap] from controller-runtime can be utilized to wrap it in a `logr` implementation.

Below is an example illustrating the use of [`zap-logfmt`][logfmt_repo] in logging.

### Example

In your `main.go` file, replace the current implementation for logs inside the `main` function:

```Go
...
// Add the zap logger flag set to the CLI. The flag set must
// be added before calling flag.Parse().
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	logf.SetLogger(logger)
...
```

With:

```Go
	import(
	...
	zaplogfmt "github.com/sykesm/zap-logfmt"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	...
)
	configLog := uzap.NewProductionEncoderConfig()
	configLog.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339Nano))
	}
	logfmtEncoder := zaplogfmt.NewEncoder(configLog)

	// Construct a new logr.logger.
	logger := zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout), zap.Encoder(logfmtEncoder))
	logf.SetLogger(logger)
```

**NOTE**: For this example, you will need to add the module `"github.com/sykesm/zap-logfmt"` to your project. Run `go get -u github.com/sykesm/zap-logfmt`.

#### Output using custom zap logger

```console
$ go run main.go
ts=2020-04-30T20:35:59.551268Z level=info logger=global msg="Printing at INFO level"
ts=2020-04-30T20:35:59.551314Z level=debug logger=global msg="Printing at DEBUG level"
ts=2020-04-30T20:35:59.551318Z level=info logger=scoped msg="Printing at INFO level"
ts=2020-04-30T20:35:59.55132Z level=debug logger=scoped msg="Printing at DEBUG level"
```

By using `sigs.k8s.io/controller-runtime/pkg/log`, your logger is propagated through `controller-runtime`. Any logs produced by `controller-runtime` code will be through your logger, and therefore have the same formatting and destination.

### Setting flags when running locally

When running locally with `make run ENABLE_WEBHOOKS=false`, you can use the `ARGS` var to pass additional flags to your operator, including the zap flags. For example:

```console
$ make run ARGS="--zap-encoder=console" ENABLE_WEBHOOKS=false
```
Make sure to have your `run` target to take `ARGS` as shown below in `Makefile`.

```makefile
# Run against the configured Kubernetes cluster in ~/.kube/config
run: manifests generate fmt vet
	go run ./main.go $(ARGS)
```

### Setting flags when deploying to a cluster

When deploying your operator to a cluster you can set additional flags using an `args` array in your operator's `container` spec in the file `config/default/manager_metrics_patch.yaml` For example:

```yaml
- op: add
  path: /spec/template/spec/containers/0/args/0
  value: --zap-log-level=debug 
- op: add
  path: /spec/template/spec/containers/0/args/0
  value: --zap-encoder=console
```

## Creating a structured log statement

There are two ways to create structured logs with `logr`. You can create new loggers using `log.WithValues(keyValues)` that include `keyValues`, a list of key-value pair `interface{}`'s, in each log record. Alternatively you can include `keyValues` directly in a log statement, as all `logr` log statements take some message and `keyValues`. The signature of `logr.Error()` has an `error`-type parameter, which can be `nil`.

An example from [`memcached_controller.go`][code_memcached_controller]:

```Go
package memcached

import (
  ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)


// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Fetch the Memcached instance
	memcached := &cachev1alpha1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Memcached resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Memcached")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMemcached(memcached)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	...
}
```

Log records will look like the following (from `log.Error()` above):

```
2020-04-27T09:14:15.939-0400	ERROR	controllers.Memcached	Failed to create new Deployment	{"memcached": "default/memcached-sample", "Deployment.Namespace": "default", "Deployment.Name": "memcached-sample"}
```

## Non-default logging

If you do not want to use `logr` as your logging tool, you can remove `logr`-specific statements without issue from your operator's code, including the `logr` setup code in `main.go`, and add your own. Note that removing `logr` setup code will prevent `controller-runtime` from logging.


[godoc_logr]:https://pkg.go.dev/github.com/go-logr/logr
[repo_zapr]:https://pkg.go.dev/github.com/go-logr/zapr
[godoc_logr_logger]:https://pkg.go.dev/github.com/go-logr/logr#Logger
[code_memcached_controller]: https://github.com/graphitehealth/operator-sdk/blob/v1.2.0/testdata/go/memcached-operator/controllers/memcached_controller.go
[logfmt_repo]:https://github.com/jsternberg/zap-logfmt
[controller_runtime_zap]:https://github.com/kubernetes-sigs/controller-runtime/tree/master/pkg/log/zap
[logging_godocs]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/log/zap#Options.BindFlags
