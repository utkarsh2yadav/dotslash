import React from "react";
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";


export default function Editor() {
    return <AceEditor 
    enableBasicAutocompletion="true"
    enableLiveAutocompletion="true"
    width="100%"
    height="90vh"
    mode="java"
    theme="monokai"
    name="UNIQUE_ID_OF_DIV"
    editorProps={{ $blockScrolling: true }}
  />
  }