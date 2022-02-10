package gremlingo

// // The processor converts the op args such that it can be serialized.
// type processor interface {
// 	convertArgs() map[string]interface{}
// }
//
// type standardProcessor struct {
// }
//
// func (sp *standardProcessor) convertArgs()
//
// type traversalProcessor struct {
// }
//
// func getProcessor(processorType string) *processor {
// 	if processorType == bytecodeProcessor {
// 		return &traversalProcessor{}
// 	} else if processorType == stringProcessor {
// 		return &standardProcessor{}
// 	}
// 	// TODO: Remote transaction session
// 	return &processor{}
// }
//