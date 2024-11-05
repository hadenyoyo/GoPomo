<div align="center">
    <h1> GoPomo </h1>
    <h3>Simple pomodoro timer made in Go.</h3>
</div>

---

Made as a simple introductory Go project.

## Usage
`./run -h` to see flags.

`-w or -work <float>`
Sets the work time in minutes, shorthand (default 25)

`-b or -break <float>`
Set break time in minutes (default 5)

`-lb or -longbreak <float>`
Explicitly sets the long break time in minutes (default 2 * break time)

`-c or -confirm`
Sets the confirm before starting next phase option (default false)

`-q or -quiet`
Disables notifications on phase completion (default false)

`-l or -loops <int>`
Sets the number of of work sessions before long break should start (default 3)

### Usage Examples
`./run -work 25 -break 5 -confirm`

`./run -w 45 -b 7.5 -lb 22 -l 5 -c`