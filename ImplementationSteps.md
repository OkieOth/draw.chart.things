# TL;DR;

This project is model driven developed. The reason is to challenge my love for models
and to practice the approach.

Find below the general steps of the implementation process

# Steps

1. Define the model that describes the input data [...](./configs/models/boxes.json)
2. Include code generation and adjust the needed types template
3. Implement and test the basis serialization/deserialization
4. Prepare the structure for the boxes implementation
5. Come up with a model to describe the printed boxes chart [...](./configs/models/boxes_document.json)
6. Generate final models
7. Implement the initial model transformation from the input to the processing model
8. Implement the processing to create the final document description
9. Render the Image
10. Save the image as file 