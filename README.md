# launch-editor

Open file with line numbers in editor from Go.


## Usage

``` go
package main

import (
	"log"

	. "github.com/Binbiubiubiu/launch-editor"
)

func main() {
	err := LaunchEditor("guess.go:59:20", "code")
	if err != nil {
		log.Fatalln(err)
	}
}
```


### Supported editors

| Value | Editor | Linux | Windows | OSX |
|--------|------|:------:|:------:|:------:|
| `appcode` | [AppCode](https://www.jetbrains.com/objc/) |  |  |✓|
| `atom` | [Atom](https://atom.io/) |✓|✓|✓|
| `atom-beta` | [Atom Beta](https://atom.io/beta) |  |  |✓|
| `brackets` | [Brackets](http://brackets.io/) |✓|✓|✓|
| `clion` | [Clion](https://www.jetbrains.com/clion/) |  |✓|✓|
| `code` | [Visual Studio Code](https://code.visualstudio.com/) |✓|✓|✓|
| `code-insiders` | [Visual Studio Code Insiders](https://code.visualstudio.com/insiders/) |✓|✓|✓|
| `codium` | [VSCodium](https://github.com/VSCodium/vscodium) |✓|✓|✓|
| `emacs` | [Emacs](https://www.gnu.org/software/emacs/) |✓| | |
| `idea` | [IDEA](https://www.jetbrains.com/idea/) |✓|✓|✓|
| `notepad++` | [Notepad++](https://notepad-plus-plus.org/download/v7.5.4.html) | |✓| |
| `pycharm` | [PyCharm](https://www.jetbrains.com/pycharm/) |✓|✓|✓|
| `phpstorm` | [PhpStorm](https://www.jetbrains.com/phpstorm/) |✓|✓|✓|
| `rubymine` | [RubyMine](https://www.jetbrains.com/ruby/) |✓|✓|✓|
| `sublime` | [Sublime Text](https://www.sublimetext.com/) |✓|✓|✓|
| `vim` | [Vim](http://www.vim.org/) |✓| | |
| `visualstudio` | [Visual Studio](https://www.visualstudio.com/vs/) | | |✓|
| `webstorm` | [WebStorm](https://www.jetbrains.com/webstorm/) |✓|✓|✓|

### Custom editor support

You can use the `LAUNCH_EDITOR` environment variable 

#### to force a specific supported editor 

```bash
LAUNCH_EDITOR=codium
```

#### to run a custom launch script

```bash
LAUNCH_EDITOR=my-editor-launcher.sh
```

```shell
# gets called with 3 args: filename, line, column
filename=$1
line=$2
column=$3

# call your editor with whatever args it expects
my-editor -l $line -c $column -f $filename
```


