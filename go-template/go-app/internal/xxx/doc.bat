rem --exclude frontend
swag fmt
swag init --output ../../api/openapi/xxx   --parseDependency --parseInternal -g doc.go
