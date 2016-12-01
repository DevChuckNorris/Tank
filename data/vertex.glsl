#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform vec3 light;

in vec3 vert;
in vec2 vertTexCoord;
in vec3 vertNormal;

out vec2 fragTexCoord;
out vec3 fragColor;

out vec3 fragNormal;
out vec3 Position_worldspace;
out vec3 Normal_cameraspace;
out vec3 EyeDirection_cameraspace;
out vec3 LightDirection_cameraspace;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vert;
    gl_Position = projection * camera * model * vec4(vert, 1);

    Position_worldspace = (model * vec4(vert, 1)).xyz;
    Normal_cameraspace = ( camera * model * vec4(vertNormal,0)).xyz;

    vec3 vertexPosition_cameraspace = ( camera * model * vec4(vert,1)).xyz;
	EyeDirection_cameraspace = vec3(0,0,0) - vertexPosition_cameraspace;

    vec3 LightPosition_cameraspace = ( camera * vec4(light,1)).xyz;
	LightDirection_cameraspace = LightPosition_cameraspace + EyeDirection_cameraspace;
}
