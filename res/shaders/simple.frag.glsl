#version 410

in vec3 outColor;

out vec4 fragColor;

uniform vec3 ambientLightColor = vec3(1, 0, 0);
uniform float ambientStrength = 1;

void main()
{
    vec3 objColor = vec3(1, 1, 1);
    fragColor = vec4(objColor * ambientLightColor * ambientStrength, 1.0);
} 
