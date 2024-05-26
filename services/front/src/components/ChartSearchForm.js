import { useForm, Controller } from "react-hook-form";
import {
  Box,
  Button,
  Checkbox,
  FormControlLabel,
  FormGroup,
} from "@mui/material";

export default function ChartSearchForm(props) {
  const { handleSubmit, control } = useForm();

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
      <Box display={"flex"} justifyContent={"flex-end"}>
        <Button type="submit" color="primary">
          Search
        </Button>
      </Box>
    </Box>
  );
}
