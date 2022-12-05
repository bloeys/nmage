//shader:vertex
#version 410

layout(location=0) in vec3 vertPosIn;
layout(location=1) in vec3 vertNormalIn;
layout(location=2) in vec2 vertUV0In;
layout(location=3) in vec3 vertColorIn;

out vec3 vertNormal;
out vec2 vertUV0;
out vec3 vertColor;
out vec3 fragPos;

//MVP = Model View Projection
uniform mat4 modelMat;
uniform mat4 viewMat;
uniform mat4 projMat;

void main()
{
    vertNormal = mat3(transpose(inverse(modelMat))) * vertNormalIn;
    vertUV0 = vertUV0In;
    vertColor = vertColorIn;
    fragPos = vec3(modelMat * vec4(vertPosIn, 1.0));

    gl_Position = projMat * viewMat * modelMat * vec4(vertPosIn, 1.0);
}

//shader:fragment
#version 410

uniform float ambientStrength = 0;
uniform vec3 ambientLightColor = vec3(1, 1, 1);

uniform vec3 lightPos1;
uniform vec3 lightColor1;

uniform sampler2D diffTex;

in vec3 vertColor;
in vec3 vertNormal;
in vec2 vertUV0;
in vec3 fragPos;

out vec4 fragColor;

void main()
{
    vec3 lightDir = normalize(lightPos1 - fragPos);  
    float diffStrength = max(0.0, dot(normalize(vertNormal), lightDir));

    vec3 finalAmbientColor = ambientLightColor * ambientStrength;
    vec4 texColor = texture(diffTex, vertUV0);
    fragColor = vec4(texColor.rgb * vertColor * (finalAmbientColor + diffStrength*lightColor1) , texColor.a);
} 
