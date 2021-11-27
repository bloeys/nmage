#version 410

in vec3 vertColor;
in vec3 vertNormal;
in vec3 fragPos;

out vec4 fragColor;

uniform float ambientStrength = 1;
uniform vec3 ambientLightColor = vec3(1, 1, 1);

uniform vec3 lightPos1;
uniform vec3 lightColor1;

void main()
{
    // vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos1 - fragPos);  
    float diffStrength = max(0.0, dot(normalize(vertNormal), lightDir));

    vec3 finalAmbientColor = ambientLightColor * ambientStrength;
    fragColor = vec4(vertColor * (finalAmbientColor + diffStrength*lightColor1), 1.0);
} 
