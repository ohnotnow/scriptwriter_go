package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

// ScriptResponse is equivalent to your Pydantic model.
type ScriptResponse struct {
	Filename string `json:"filename" jsonschema_description:"The filename for the generated script"`
	Content  string `json:"content" jsonschema_description:"The content of the generated script"`
}

// GenerateSchema creates a JSON schema for a given generic type.
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

// Generate the JSON schema at initialization time.
var ScriptResponseSchema = GenerateSchema[ScriptResponse]()

// buildPrompt constructs the prompt using the provided requirements.
func buildPrompt(requirements string) string {
	return fmt.Sprintf(`
You are a commandline script writer.
You are given a set of requirements for a script.
You need to write a script that satisfies the requirements.
You can choose to write the script in either Python, NodeJS, PHP or BASH.
The script will primarily be used on MacOS or Linux - please remember to account for the differences between the GNU tooling and BSD/Darwin tooling.
You need to return the filename, and the content of the script.

<user-requirements>
%s
</user-requirements>
`, requirements)
}

func main() {
	// Read requirements from the user.
	fmt.Print("Enter the requirements for the script: ")
	reader := bufio.NewReader(os.Stdin)
	requirements, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	requirements = strings.TrimSpace(requirements)
	prompt := buildPrompt(requirements)

	// Initialize the OpenAI client.
	client := openai.NewClient()
	ctx := context.Background()

	// Define the JSON schema parameter for structured outputs.
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("script"),
		Description: openai.F("Script generation based on requirements"),
		Schema:      openai.F(ScriptResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Query the Chat Completions API using the structured outputs feature.
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		// Use a model that supports structured outputs.
		Model: openai.F(openai.ChatModelO3Mini),
	})
	if err != nil {
		panic(err.Error())
	}

	// Unmarshal the structured JSON response into a ScriptResponse struct.
	var scriptResp ScriptResponse
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &scriptResp)
	if err != nil {
		panic(err.Error())
	}

	// Determine the output filename, with a fallback.
	outputFilename := scriptResp.Filename
	if outputFilename == "" {
		outputFilename = "output.sh"
	}

	// Write the generated script to a file and make it executable.
	err = ioutil.WriteFile(outputFilename, []byte(scriptResp.Content), 0755)
	if err != nil {
		panic(err.Error())
	}

	// Print the content and file path.
	fmt.Println(scriptResp.Content)
	fmt.Printf("### Script written to %s\n", outputFilename)
}
