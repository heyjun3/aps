import TextField from "@mui/material/TextField";
import Autocomplete from "@mui/material/Autocomplete";

export function PageBox(props) {
  const { setLimit } = props
  return (
    <div>
      <Autocomplete
        disablePortal
        id="count-box"
        options={counts}
        sx={{
          width: 110,
          paddingTop: 2,
          marginRight: 2,
          marginLeft: "auto",
        }}
        renderInput={(params) => <TextField {...params} label="Count" />}
        onChange={(_, newValue) => setLimit(newValue)}
        defaultValue={100}
      />
    </div>
  );
}

const counts = ["100", "150"];
