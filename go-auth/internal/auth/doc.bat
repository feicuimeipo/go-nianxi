rem --exclude frontend
swag fmt
swag init --output ../../api/openapi/auth   --parseDependency --parseInternal -g doc.go
