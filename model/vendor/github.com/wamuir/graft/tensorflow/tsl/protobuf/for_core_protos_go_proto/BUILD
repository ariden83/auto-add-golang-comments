load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "for_core_protos_go_proto",
    srcs = [
        "bfc_memory_map.pb.go",
        "distributed_runtime_payloads.pb.go",
        "error_codes.pb.go",
        "histogram.pb.go",
        "rpc_options.pb.go",
        "status.pb.go",
        "test_log.pb.go",
    ],
    importpath = "github.com/wamuir/graft/tensorflow/tsl/protobuf/for_core_protos_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/wrapperspb",
    ],
)

alias(
    name = "go_default_library",
    actual = ":for_core_protos_go_proto",
    visibility = ["//visibility:public"],
)
