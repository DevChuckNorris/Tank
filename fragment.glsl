#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 fragColor;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fract(fragTexCoord));
}
