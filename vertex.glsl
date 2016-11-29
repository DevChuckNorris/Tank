#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;
out vec3 fragColor;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vert;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
