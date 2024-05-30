import { useForm, Controller } from "react-hook-form";
import {
  Box,
  Button,
  Checkbox,
  FormControlLabel,
  FormGroup,
  TextField,
} from "@mui/material";

export default function ChartSearchForm(props) {
  const { handleSubmit, control } = useForm();
  const excludeKeywords = props.excludeKeywords?.split(",").join(" ");
  return (
    <Box
      component="form"
      width={"60%"}
      margin={"auto"}
      onSubmit={handleSubmit(props.onSubmit)}
    >
      <Controller
        name="rankLine"
        control={control}
        defaultValue={false}
        render={({ field, formState: { errors } }) => (
          <FormGroup {...field}>
            <FormControlLabel control={<Checkbox />} label="Rank Line" />
          </FormGroup>
        )}
      />
      <Controller
        name="excludeKeywords"
        control={control}
        defaultValue={excludeKeywords}
        render={({ field, formState: { errors } }) => (
          <TextField
            {...field}
            id="standard-basic"
            variant="standard"
            label="Exclude Keywords"
            defaultValue={excludeKeywords}
            fullWidth
          />
        )}
      />
      <Box display={"flex"} justifyContent={"flex-end"}>
        <Button type="submit" color="primary">
          Search
        </Button>
      </Box>
    </Box>
  );
}
