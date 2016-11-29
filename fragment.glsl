#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 fragColor;

out vec4 outputColor;

void main() {
    outputColor = vec4(fragColor, 1.0); // texture(tex, fragTexCoord);
}
