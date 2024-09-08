import { createTheme } from "@mui/material";

export const createMuiTheme = (tgTheme) => {
    return createTheme({
        palette: {
            common :{
                black: 'red'
            },
            text: {
                primary: tgTheme.textColor,
                secondary: tgTheme.accentTextColor
            },
            background: {
                default: tgTheme.bgColor,
                paper: tgTheme.secondaryBgColor
            },
            primary: {
                main: tgTheme.accentTextColor
            },
            secondary: {
                main: tgTheme.textColor
            },
            action: {
                active: tgTheme.accentTextColor
            }

        },

          
    });
}