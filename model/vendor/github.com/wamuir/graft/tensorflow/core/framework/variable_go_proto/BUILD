load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "variable_go_proto",
    srcs = ["variable.pb.go"],
    importpath = "github.com/wamuir/graft/tensorflow/core/framework/variable_go_proto",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)

alias(
    name = "go_default_library",
    actual = ":variable_go_proto",
    visibility = ["//visibility:public"],
)
