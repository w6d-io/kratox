#!/bin/bash

ROOT=$PWD

walk() {
  let LOL=0
  for d in * ; do
    if [ -d "$d" -a ${d##*/} != "vendor" ]; then
      (cd -- "$d" && walk)
      continue
    fi
      count=`ls -1 *.go 2>/dev/null | wc -l`
      if [ $count != 0 -a $LOL == 0 ]; then
        let LOL++
        echo "generate README.md for this folder $PWD"
        goreadme --variabless -skip-examples -constants -credit=false -methods -functions -factories -recursive -types > README.md
        printf -v path "%s/%s" $PWD $d
        sed -i "s,/$d,${path#"$ROOT"},g" README.md
        continue
      fi
      if [ $count != 0 ]; then
        printf -v path "%s/%s" $PWD $d
        sed -i "s,/$d,${path#"$ROOT"},g" README.md
      fi
  done
}

read -p "If vendor directory exist, it will be deleted. Confirm? [yY/n] " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "processing..."
  find . -name "vendor" -ok rm -rf {} \;
  walk
fi
