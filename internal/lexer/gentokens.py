import sys

PACKAGE = "lexer"

file = open(sys.argv[1])
output = open(sys.argv[2], "w")

states = set()

output.write(f"package {PACKAGE}\n\n")

output.write(f"import \"regexp\"\n\n")

output.write("var tokens = []token{\n")

for line in file:
    try:
        state, name, expr = line.strip().split(' ', 2)
    except:
        continue
    state = state.strip()
    name = name.strip()
    expr = expr.strip()

    push = 'stateKeep()'
    skip = 'false'
    if '->' in name:
        name, push = name.split('->', 1)
        push = f'statePush(state{push})'
    elif name.endswith('<-'):
        push = 'statePop()'
        name = name[:-2]
    
    if name[0] == '.':
        skip = 'true'
        name = name[1:]

    
    states.add(state)

    expr = expr.replace("\\", "\\\\").replace("\"", "\\\"")

    output.write(f'    {{ state: state{state}, name: "{name}", stateChange: {push}, skip: {skip}, expr: regexp.MustCompile("^({expr})") }},\n')

output.write("}\n\n")

output.write("const (\n")
for i, state in enumerate(sorted(states)):
    if i == 0:
        output.write(f"    state{state} state = iota\n")
    else:
        output.write(f"    state{state}\n")
output.write(")\n")

file.close()
output.close()