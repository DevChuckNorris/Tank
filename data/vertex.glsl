#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

layout(location = 0) in vec3 vert;
layout(location = 1) in vec2 vertTexCoord;
layout(location = 2) in vec3 vertNormal;

out vec3 fragPos;
out vec3 fragNormal;
out vec2 fragTexCoord;

void main() {
    gl_Position = projection * camera * model * vec4(vert, 1);
    fragPos = vec3(model * vec4(vert, 1));
    fragNormal = mat3(transpose(inverse(model))) * vertNormal;
    fragTexCoord = vertTexCoord;
}
