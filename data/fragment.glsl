#version 330

uniform sampler2D tex;
uniform vec3 light;

in vec3 fragPos;
in vec3 fragNormal;
in vec2 fragTexCoord;

out vec3 outputColor;

void main() {
    /*float ambientStrength = 0.1f;
    vec3 ambient = ambientStrength * vec3(1, 1, 1);

    vec3 norm = normalize(fragNormal);
    vec3 lightDir = normalize(light - fragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * vec3(1, 1, 1);

    vec3 result = (ambient + diffuse) * vec3(texture(tex, fragTexCoord));
    outputColor = result;*/


    //float fDiffuseIntensity = max(0.0, dot(normalize(vNormal), -sunLight.vDirection));
    vec3 texColor = vec3(texture(tex, fragTexCoord));

    float ambientIntensity = 0.25f;
    float diffuseIntensity = max(0.0, dot(normalize(fragNormal), -light));
    outputColor = texColor * (ambientIntensity + diffuseIntensity);
}
