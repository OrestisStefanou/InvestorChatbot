package services

type RagResponsesStore interface {
	StoreRagResponse(
		ragTopic Topic,
		prompt string,
		conversation []Message,
		response string,
	) error
}

type BaseRag struct {
	topic         Topic
	llm           Llm
	responseStore RagResponsesStore
}

// GenerateLllmResponse streams an LLM-generated response for the given prompt and conversation.
//
// This method performs the following steps:
//  1. Prepends the user prompt to the provided conversation history.
//  2. Asynchronously calls the underlying LLM to generate a response in chunks.
//  3. Sends each response chunk to the provided responseChannel as it becomes available.
//  4. Accumulates all chunks into a complete response message.
//  5. Stores the full conversation and response in the configured RagResponsesStore.
//
// Parameters:
//
//	prompt           - The user prompt to send to the LLM.
//	conversation     - The conversation history prior to this request.
//	responseChannel  - A channel to stream partial LLM response chunks back to the caller.
//
// Returns:
//
//	error - Any error encountered during response generation or storage.
//
// Notes:
//   - The responseChannel is closed once the LLM finishes generating the response.
//   - If the LLM returns an error mid-stream, generation stops and the error is returned.
//   - The complete response (not just streamed chunks) is persisted via the RagResponsesStore.
func (r *BaseRag) GenerateLllmResponse(
	prompt string,
	conversation []Message,
	responseChannel chan<- string,
) error {
	conversation = append([]Message{{Content: prompt, Role: User}}, conversation...)

	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	go func() {
		if err := r.llm.GenerateResponse(conversation, chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				shouldExit = true
				close(responseChannel)
				continue
			}
			responseMessage += chunk
			responseChannel <- chunk
		case err := <-errorChannel:
			if err != nil {
				return err
			}
		}
	}

	return r.responseStore.StoreRagResponse(
		r.topic,
		prompt,
		conversation,
		responseMessage,
	)
}

// streamChunks handles the common logic for streaming text responses from a generator function
// that produces output in chunks. It takes care of:
//
//  1. Creating the chunk and error channels.
//  2. Running the provided generate() function in a goroutine.
//  3. Forwarding each chunk to the given responseChannel.
//  4. Accumulating all chunks into a single final string.
//  5. Returning the final accumulated string or an error.
//
// Parameters:
//   - generate: A function that accepts a `chan<- string` and sends output chunks to it.
//     It should close the provided channel when finished. If it returns an error,
//     the streaming process will stop and the error will be returned.
//   - responseChannel: A channel where streamed chunks will be sent for immediate consumption.
//
// Returns:
//   - A string containing the concatenated result of all chunks.
//   - An error if one occurred during generation.
//
// Usage example:
//
//	finalMessage, err := streamChunks(
//	    func(ch chan<- string) error {
//	        return llm.GenerateResponse(conversation, ch)
//	    },
//	    responseChannel,
//	)
//
// Notes:
//   - responseChannel is closed when the stream ends.
//   - generate() is responsible for closing its own output channel.
func streamChunks(
	generate func(chan<- string) error,
	responseChannel chan<- string,
) (string, error) {
	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	// Run the generator in a goroutine
	go func() {
		if err := generate(chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				shouldExit = true
				close(responseChannel)
				continue
			}
			responseMessage += chunk
			responseChannel <- chunk
		case err := <-errorChannel:
			if err != nil {
				return "", err
			}
		}
	}
	return responseMessage, nil
}
