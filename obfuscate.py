import random 
import string

# Global Variables
random_values = []

# Generate Random String of 5-8 Characters
def generate_random_string():
    global random_values
    while [ True ]:
        rand = "".join(random.choice(string.ascii_letters) for _ in range(random.randint(5,9)))
        if rand not in random_values:
            random_values.append(rand)
            return rand

# Character --> String Map
def map_characters(cmd, space_token, equals_token):
    mapping = {}
    mapping[" "] = space_token
    mapping["="] = equals_token
    for part in cmd:
        if part not in mapping:
            mapping[part] = generate_random_string()
    return mapping

# Payload Object
class Create_Payload:
    def __init__(self, cmd):
        self.cmd = cmd
        self.bad_char = string.punctuation
        self.set_token = generate_random_string()
        self.space_token = generate_random_string()
        self.equals_token = generate_random_string()
        self.mapping = map_characters(self.cmd, self.space_token, self.equals_token)

    # Initial Setup for Obfuscation
    def initial_setup(self):
        initial = [
            f"set {self.set_token}=set",
            f"%{self.set_token}% {self.space_token}= ",
            f"%{self.set_token}%%{self.space_token}%{self.equals_token}=="
        ]
        return initial

    # Use Character Mapping to Craft Payload (After Initial Setup)
    def obfuscate(self, initial_payload):
        new_payload = []
        command_parts = ""
        for part in self.cmd:
            if part in self.bad_char:
                command_parts += part                
            else:
                obfuscated_part = self.mapping[part]
                command_parts += f"%{obfuscated_part}%"
                new_payload += [
                    f"%{self.set_token}%%{self.space_token}%{obfuscated_part}%{self.equals_token}%{part}"
                ]
        random.shuffle(new_payload)
        return "{}\n{}{}".format("\n".join(initial_payload), "&".join(new_payload), "\n{}".format(command_parts))

def main():
    # Create Payload Object w/ Command to Obfuscate
    command = input("Command> ").strip()
    payload = Create_Payload(command)
    initial_payload = payload.initial_setup()
    obfuscated_payload = payload.obfuscate(initial_payload)

    # Write to .bat File
    print("\n[+] Command: {}\n[+] Payload Size: {} Characters\n[+] Payload:\n{}\n\n[+] Writing to payload.bat...".format(command, len(obfuscated_payload), obfuscated_payload))
    with open("payload.bat", "w") as file_handle:
        try:
            file_handle.write(obfuscated_payload)
            print("[+] Payload Successfully Written!")
        except:
            print("[-] Unknown File Error has Occured.")

if __name__ == "__main__":
    main()
