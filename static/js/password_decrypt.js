
async function decryptPassword(encryptedContent, privateKeyArmored, passphrase) {
    try {
        const privateKey = await openpgp.readPrivateKey({ armoredKey: privateKeyArmored });
        
        const decryptedKey = await openpgp.decryptKey({
            privateKey,
            passphrase
        });

        console.log("Reading message");
        const encryptedMessage = await openpgp.readMessage({
            binaryMessage: encryptedContent
        });

        const decrypted = await openpgp.decrypt({
            message: encryptedMessage,
            decryptionKeys: decryptedKey
        });

        return decrypted.data;
    } catch (error) {
        console.error('Decryption failed:', error);
        throw error;
    }
}

window.passphrase = "";
async function handlePasswordDecrypt() {
    const passwordContent = document.getElementById('password-content');
    const encryptedContent = passwordContent.querySelector('pre').textContent;

    const binaryString = atob(encryptedContent);
    const len = binaryString.length;
    const uint8Array = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        uint8Array[i] = binaryString.charCodeAt(i);
    }

    const privateKey = document.getElementById('privateKey').value;

    // if (window.passphrase?.length == 0){
    //     let passInput = document.getElementById("passphrase")
    //     window.passphrase = passInput?.value
    //     setTimeout(() => {window.passphrase = ""; 
    //         if (passInput){
    //             passInput.style.display = "block"
    //         }
    //
    //     }, 1000 * 60 * 15);
    // }

    // if (!window.passphrase) {
    //     let passInput = document.getElementById("passphrase")
    //     window.passphrase = passInput?.value
    //     setTimeout(() => {window.passphrase = ""; 
    //         if (passInput){
    //             passInput.style.display = "block"
    //         }
    //
    //     }, 1000 * 60);
    // }

    try {
        let password;
        let password_hide = false;
        if (!window.passphrase) {
            let passInput = document.getElementById("passphrase")
            password_hide = true
            password = passInput.value
        } else {
            password = window.passphrase
        }
        const decryptedContent = await decryptPassword(uint8Array, privateKey, password);
        passwordContent.querySelector('pre').textContent = decryptedContent;
        document.getElementById("passphrase").style.display = "none"

        if (password_hide){
            let passInput = document.getElementById("passphrase")
            window.passphrase = passInput?.value
            setTimeout(() => {window.passphrase = ""; 
                if (passInput){
                    passInput.style.display = "block"
                }

            }, 1000 * 60);
        }
    } catch (error) {
        alert('Failed to decrypt password: ' + error.message);
    }
}

// document.addEventListener('htmx:afterSwap', function(evt) {
//     if (evt.detail.tarkget.id === 'password-content') {
//         const passwordMenu = document.getElementById('passwordMenu');
//         passwordMenu.style.display = 'block';
//     }
// }); 
