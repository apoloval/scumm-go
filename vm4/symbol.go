package vm4

import "github.com/apoloval/scumm-go/vm"

// DefaultSymbolTable returns a symbol table with the default variables used in SCUMM v4.
func DefaultSymbolTable() *vm.SymbolTable {
	st := vm.NewSymbolTable()
	st.Declare(vm.SymbolTypeLabel, "START", 0)
	for i, name := range defaultVarNames {
		if name != "" {
			st.Declare(vm.SymbolTypeVar, name, uint16(i))
		}
	}
	return st
}

var defaultVarNames = []string{
	/* 0 */
	"VAR_RESULT",
	"VAR_EGO",
	"VAR_CAMERA_POS_X",
	"VAR_HAVE_MSG",
	/* 4 */
	"VAR_ROOM",
	"VAR_OVERRIDE",
	"VAR_MACHINE_SPEED",
	"VAR_ME",
	/* 8 */
	"VAR_NUM_ACTOR",
	"VAR_CURRENT_LIGHTS",
	"VAR_CURRENTDRIVE",
	"VAR_TMR_1",
	/* 12 */
	"VAR_TMR_2",
	"VAR_TMR_3",
	"VAR_MUSIC_TIMER",
	"VAR_ACTOR_RANGE_MIN",
	/* 16 */
	"VAR_ACTOR_RANGE_MAX",
	"VAR_CAMERA_MIN_X",
	"VAR_CAMERA_MAX_X",
	"VAR_TIMER_NEXT",
	/* 20 */
	"VAR_VIRT_MOUSE_X",
	"VAR_VIRT_MOUSE_Y",
	"VAR_ROOM_RESOURCE",
	"VAR_LAST_SOUND",
	/* 24 */
	"VAR_CUTSCENEEXIT_KEY",
	"VAR_TALK_ACTOR",
	"VAR_CAMERA_FAST_X",
	"VAR_SCROLL_SCRIPT",
	/* 28 */
	"VAR_ENTRY_SCRIPT",
	"VAR_ENTRY_SCRIPT2",
	"VAR_EXIT_SCRIPT",
	"VAR_EXIT_SCRIPT2",
	/* 32 */
	"VAR_VERB_SCRIPT",
	"VAR_SENTENCE_SCRIPT",
	"VAR_INVENTORY_SCRIPT",
	"VAR_CUTSCENE_START_SCRIPT",
	/* 36 */
	"VAR_CUTSCENE_END_SCRIPT",
	"VAR_CHARINC",
	"VAR_WALKTO_OBJ",
	"VAR_DEBUGMODE",
	/* 40 */
	"VAR_HEAPSPACE",
	"",
	"VAR_RESTART_KEY",
	"VAR_PAUSE_KEY",
	/* 44 */
	"VAR_MOUSE_X",
	"VAR_MOUSE_Y",
	"VAR_TIMER",
	"VAR_TIMER_TOTAL",
	/* 48 */
	"VAR_SOUNDCARD",
	"VAR_VIDEOMODE",
	"VAR_MAINMENU_KEY",
	"VAR_FIXEDDISK",
	/* 52 */
	"VAR_CURSORSTATE",
	"VAR_USERPUT",
	"VAR_V5_TALK_STRING_Y",
}
