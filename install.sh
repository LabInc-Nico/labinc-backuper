#!/bin/sh

set -e

if [ -z "${BIN_DIR}" ]; then
	BIN_DIR=$(pwd)
fi

THE_ARCH_BIN=""
DEST=${BIN_DIR}/labinc-backuper

OS=$(uname -s)
ARCH=$(uname -m)

if type "tput" >/dev/null 2>&1; then
	bold=$(tput bold || true)
	italic=$(tput sitm || true)
	normal=$(tput sgr0 || true)
fi

case ${OS} in
Linux*)
	case ${ARCH} in
	aarch64)
		THE_ARCH_BIN=""
		;;
	x86_64)
		THE_ARCH_BIN="labinc-backuper_$(uname)_$(uname -m).tar.gz"
		;;
	*)
		THE_ARCH_BIN=""
		;;
	esac
	;;
Darwin*)
	case ${ARCH} in
	arm64)
		THE_ARCH_BIN=""
		;;
	*)
		THE_ARCH_BIN=""
		;;
	esac
	;;
Windows | MINGW64_NT*)
	echo "❗ Use WSL to run labinc-backuper on Windows: https://learn.microsoft.com/windows/wsl/"
	exit 1
	;;
*)
	THE_ARCH_BIN=""
	;;
esac

if [ -z "${THE_ARCH_BIN}" ]; then
	echo "❗ labinc-backuper is not supported on ${OS} and ${ARCH}"
	exit 1
fi


echo "📦 Downloading ${bold}labinc-backuper${normal} for ${OS} (${ARCH}):"

# check if $DEST is writable and suppress an error message
touch "${DEST}" 2>/dev/null

if type "curl" >/dev/null 2>&1; then
	curl -L --progress-bar "https://github.com/labinc-nico/labinc-backuper/releases/latest/download/${THE_ARCH_BIN}" -o "${DEST}"
elif type "wget" >/dev/null 2>&1; then
	wget "https://github.com/labinc-nico/labinc-backuper/releases/latest/download/${THE_ARCH_BIN}" -O "${DEST}"
else
	echo "❗ Please install ${italic}curl${normal} or ${italic}wget${normal} to download labinc-backuper"
	exit 1
fi

tar -zxf "${DEST}"
chmod u+x labinc-backuper
mv .labinc-backuper.example.yaml .labinc-backuper.yaml

echo
echo "🥳 labinc-backuper downloaded successfully to ${italic}${DEST}${normal}"
echo "👉 edit the config file: ${bold}.labinc-backuper.yaml${normal}"
echo "🔧 Run it with ${bold}./labinc-backuper${normal}"
echo
echo "⭐ If you like labinc-backuper, please give it a star on GitHub: ${italic}https://github.com/labinc-nico/labinc-backuper${normal}"