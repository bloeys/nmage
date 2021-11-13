#version 460

in vec3 vertPos;
in vec3 vertColor;

out vec3 outColor;

//MVP = Model View Projection
uniform mat4 modelMat;
uniform mat4 viewMat;
uniform mat4 projMat;

void main()
{
    outColor = vertColor;
    gl_Position = projMat * viewMat * modelMat * (vec4(vertPos, 1.0)); // vec4(vertPos.x, vertPos.y, vertPos.z, 1.0)
}