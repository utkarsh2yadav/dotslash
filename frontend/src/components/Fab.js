
import Button from '@mui/material/Button';
import { PlayArrow, StopRounded } from "@mui/icons-material";
import { createTheme } from '@mui/material/styles';

const theme = createTheme();

theme.spacing(2);
export default function Faab() {
    return (
        <>
        

<Button
    style={{
        borderRadius: 22,
        backgroundColor: "#1c2566",
        padding: "10px 45px",
        fontSize: "14px"
    }}
    variant="contained"
    >
    <PlayArrow sx={{ mr: 1 }} />
    START
</Button>

</>


    );
  }

  