package spec_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec"
)

func versionsFrom(first spec.DataVersion) []spec.DataVersion {
	versions := []spec.DataVersion{}
	for v := first; v.String() != "unknown"; v++ {
		versions = append(versions, v)
	}

	return versions
}

func fill(v reflect.Value, depth int) {
	if depth == 0 {
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() && v.CanSet() && v.Type().Elem().Kind() == reflect.Struct {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if !v.IsNil() {
			fill(v.Elem(), depth-1)
		}
	case reflect.Struct:
		for i := range v.NumField() {
			if v.Type().Field(i).IsExported() {
				fill(v.Field(i), depth-1)
			}
		}
	default:
	}
}

func isUnknownVersion(out reflect.Value) bool {
	switch {
	case out.Type() == reflect.TypeFor[string]():
		return out.String() == "unknown version"
	case out.Type().Implements(reflect.TypeFor[error]()):
		if out.IsNil() {
			return false
		}

		msg := out.Interface().(error).Error()

		return strings.Contains(msg, "unknown version") || strings.Contains(msg, "unsupported version")
	default:
		return false
	}
}

func callMethod(method reflect.Value) (outs []reflect.Value, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	args := make([]reflect.Value, method.Type().NumIn())
	for i := range args {
		args[i] = reflect.Zero(method.Type().In(i))
	}

	return method.Call(args), false
}

// TestVersionedTypesHandleLaterForks fails when a versioned type that handles
// Gloas falls through to its "unknown version" branch for a later fork.
func TestVersionedTypesHandleLaterForks(t *testing.T) {
	types := []any{
		spec.VersionedAggregateAndProof{},
		spec.VersionedAttestation{},
		spec.VersionedAttesterSlashing{},
		spec.VersionedBeaconBlock{},
		spec.VersionedBeaconBlockBody{},
		spec.VersionedBeaconState{},
		spec.VersionedExecutionPayload{},
		spec.VersionedExecutionPayloadEnvelope{},
		spec.VersionedExecutionRequests{},
		spec.VersionedIndexedAttestation{},
		spec.VersionedSignedAggregateAndProof{},
		spec.VersionedSignedBeaconBlock{},
		spec.VersionedSignedExecutionPayloadBid{},
		spec.VersionedSignedExecutionPayloadEnvelope{},
	}

	for _, typ := range types {
		structType := reflect.TypeOf(typ)
		if _, ok := structType.FieldByName("Gloas"); !ok {
			continue
		}

		for _, version := range versionsFrom(spec.DataVersionGloas) {
			t.Run(structType.Name()+"/"+version.String(), func(t *testing.T) {
				value := reflect.New(structType)
				fill(value.Elem(), 8)
				value.Elem().FieldByName("Version").Set(reflect.ValueOf(version))

				for i := range value.NumMethod() {
					name := value.Type().Method(i).Name
					if strings.HasPrefix(name, "Unmarshal") {
						continue
					}

					outs, panicked := callMethod(value.Method(i))
					if panicked {
						continue
					}

					for _, out := range outs {
						if isUnknownVersion(out) {
							t.Errorf("%s.%s does not handle %s", structType.Name(), name, version)
						}
					}
				}
			})
		}
	}
}
