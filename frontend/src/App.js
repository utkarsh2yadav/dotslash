import { Grid } from '@mui/material';
import { Terminal } from 'xterm';
import './App.css';
import Editor from './components/Editor.js';
import SearchAppBar from './components/SearchAppBar.js';
import Console from './components/Console.js';


// import Term from './components/Term'


function App() {

  let xterm = new Terminal()

  return (
    <div className="App">
      <SearchAppBar />
      <Grid container style={{ height: "90vh" }} spacing={"5"}>
        <Grid item md={6}>
          <Editor xterm={xterm} />
        </Grid>
        
        <Grid item md={6}>
          <Console xterm={xterm} />
          {/* <Term xterm={xterm} /> */}
        </Grid>
      </Grid>
      
    </div>
  );
}

export default App;
