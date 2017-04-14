package traces

import (
	"io"
	"log"

	"github.com/fatih/color"
)

/*DebugTracer is a Logger object to output logs to debug the program
 *ErrorTracer is a Logger object to output logs for runtime errors
 *InfoTracer is a Logger object to output basic informations about the program
 *WarningTracer is a Logger object to output runtime warnings - warnings are informative messages that are not considered as errors
 */
var (
	DebugTracer   *log.Logger
	ErrorTracer   *log.Logger
	InfoTracer    *log.Logger
	WarningTracer *log.Logger
)

/*InitTraces is a function that initialize Loggers
 */
func InitTraces(debugHandle io.Writer, errorHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer) {

	blue := color.New(color.FgBlue).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	/*Initialize the debug field
	 */
	DebugTracer = log.New(debugHandle, blue("DEBUG: "), log.Ldate|log.Ltime|log.Lshortfile)

	/*Initialize the error field
	 */
	ErrorTracer = log.New(errorHandle, red("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)

	/*Initialize the info field
	 */
	InfoTracer = log.New(infoHandle, cyan("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)

	/*Initialize the warning field
	 */
	WarningTracer = log.New(warningHandle, yellow("WARNING: "), log.Ldate|log.Ltime|log.Lshortfile)

}
