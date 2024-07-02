rem --exclude frontend
swag fmt
swag init --output ../../api/openapi/admin   --parseDependency --parseInternal -g doc.go
