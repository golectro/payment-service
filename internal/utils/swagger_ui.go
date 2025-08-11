package utils

func SwaggerHTML() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Swagger UI</title>
		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css">
		<style>
			html { box-sizing: border-box; overflow-y: scroll; }
			*, *:before, *:after { box-sizing: inherit; }
			body { margin:0; background: #fafafa; }
			.swagger-ui .topbar { display: none; }
		</style>
	</head>
	<body>
		<div id="swagger-ui"></div>
		<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
		<script>
			window.onload = function() {
				SwaggerUIBundle({
					url: 'swagger.json',
					dom_id: '#swagger-ui',
					deepLinking: true,
					layout: "BaseLayout",
					presets: [
						SwaggerUIBundle.presets.apis,
					]
				});
			}
		</script>
	</body>
	</html>
	`
}
