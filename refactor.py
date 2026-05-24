import os
import glob
import re

def process_file(filepath):
    if os.path.basename(filepath) in ['image_test.go', 'common_test.go']:
        return

    with open(filepath, 'r') as f:
        content = f.read()

    # Skip if no expectedValue inline assignment is found
    if 'expectedValue := "' not in content:
        return

    print(f"Processing {filepath}")

    filename_base = os.path.basename(filepath).replace('_test.go', '')
    suite_name_const = f"{filename_base.upper()}_TEST_SUITE_NAME"
    suite_name_val = filename_base

    # Update imports
    imports_match = re.search(r'import \((.*?)\)', content, re.DOTALL)
    if imports_match:
        imports = imports_match.group(1)
        new_imports = imports
        if '"os"' not in new_imports:
            new_imports += '\t"os"\n'
        if '"fmt"' not in new_imports:
            new_imports += '\t"fmt"\n'
        content = content.replace(imports, new_imports)

    # Add constant if not present
    if suite_name_const not in content:
        content = content.replace(')\n', f')\n\nconst {suite_name_const} = "{suite_name_val}"\n', 1)

    # Regex to match the old block inside test functions
    # We also need to extract the test name.
    
    # We will split by "func Test", process each function
    parts = content.split("func Test")
    new_parts = [parts[0]]
    for part in parts[1:]:
        # Find function name
        test_name_match = re.match(r'([A-Za-z0-9_]+)\(.*?{', part)
        if not test_name_match:
            new_parts.append("func Test" + part)
            continue
            
        test_name = "Test" + test_name_match.group(1)
        
        # Regex to find the replacement block
        # We look for: expectedValue := "..." and optionally the following lines since they might vary slightly,
        # but the instructions say they are:
        # expectedValue := "..."
        # obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
        # expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
        # if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
        #     fmt.Println("Expected:\n", expectedValueBase64)
        #     fmt.Println("Obtained:\n", obtainedValueBase64)
        # }
        
        # A simpler approach: find expectedValue := "(long string)"
        # and replace it with the template. The following lines can be left intact if they are the same!
        # Wait, the prompt says "Replace `expectedValue := "..."` followed by `layerEntry.GetBasicAnsiStringAsBase64()`." No, it says:
        # "Identify all *_test.go files... that contain expectedValue := "..." followed by layerEntry.GetBasicAnsiStringAsBase64()." Wait, usually obtainedValue := layerEntry.GetBasicAnsiStringAsBase64() comes BEFORE expectedValue := "...".
        
        # Let's replace the single line: expectedValue := "..."
        # With:
        # if UpdateMasterImages(false, SUITE_NAME, "TestName", obtainedValue) {
        #     return
        # }
        # expectedValueBytes, err := os.ReadFile(constants.MasterImagesPath + SUITE_NAME + "/" + "TestName.base64")
        # assert.Nil(test, err, "Reading expected value file should not produce an error")
        # expectedValue := string(expectedValueBytes)
        
        replacement = f"""if UpdateMasterImages(false, {suite_name_const}, "{test_name}", obtainedValue) {{
\t\treturn
\t}}
\texpectedValueBytes, err := os.ReadFile(constants.MasterImagesPath + {suite_name_const} + "/" + "{test_name}.base64")
\tassert.Nil(test, err, "Reading expected value file should not produce an error")
\texpectedValue := string(expectedValueBytes)"""
        
        # Find expectedValue := "..."
        part_modified = re.sub(r'expectedValue := "[^"]+"', replacement, part)
        new_parts.append("func Test" + part_modified)
        
    final_content = "".join(new_parts)
    
    # Formatting might be slightly off but we can run `go fmt` later.
    with open(filepath, 'w') as f:
        f.write(final_content)

for filepath in glob.glob('*_test.go'):
    process_file(filepath)
