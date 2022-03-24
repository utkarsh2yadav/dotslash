import { Grid } from '@mui/material';
import './App.css';
import Editor from './components/Editor.js';
import SearchAppBar from './components/SearchAppBar.js';
import BasicSpeedDial from './components/SpeedDial.js';
function App() {
  return (
    <div className="App">
      <SearchAppBar />
      <Grid container spacing={0}>
        <Grid item xs={6} md={6}>
        <Editor />
        <BasicSpeedDial />

        </Grid>
        
        
      </Grid>
      
    </div>
  );
}

export default App;
