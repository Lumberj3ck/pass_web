
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
    const passphrase = prompt('Please enter your passphrase:');

    console.log(privateKey);
    if (!passphrase) {
        return;
    }

    try {
        const decryptedContent = await decryptPassword(uint8Array, privateKey, passphrase);
        passwordContent.querySelector('pre').textContent = decryptedContent;
    } catch (error) {
        alert('Failed to decrypt password: ' + error.message);
    }
}

// document.addEventListener('htmx:afterSwap', function(evt) {
//     if (evt.detail.target.id === 'password-content') {
//         const passwordMenu = document.getElementById('passwordMenu');
//         passwordMenu.style.display = 'block';
//     }
// }); 