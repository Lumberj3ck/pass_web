async function deriveKey(password, salt) {
    const enc = new TextEncoder();
    const keyMaterial = await window.crypto.subtle.importKey(
        'raw',
        enc.encode(password),
        { name: 'PBKDF2' },
        false,
        ['deriveKey']
    );
    return window.crypto.subtle.deriveKey(
        {
            name: 'PBKDF2',
            salt: salt,
            iterations: 100000,
            hash: 'SHA-256'
        },
        keyMaterial,
        { name: 'AES-GCM', length: 256 },
        true,
        ['encrypt', 'decrypt']
    );
}

async function encryptData(key, data, privateKeyPassword) {
    const enc = new TextEncoder();
    const iv = window.crypto.getRandomValues(new Uint8Array(12));

    const dataObj = { "privateKey": data, "privateKeyPassword": privateKeyPassword };
    const jsonString = JSON.stringify(dataObj);

    const encryptedContent = await window.crypto.subtle.encrypt(
        {
            name: 'AES-GCM',
            iv: iv
        },
        key,
        enc.encode(jsonString)
    );
    return {
        iv: iv,
        data: new Uint8Array(encryptedContent)
    };
}

async function decryptData(key, encryptedData) {
    const dec = new TextDecoder();
    const decryptedContent = await window.crypto.subtle.decrypt(
        {
            name: 'AES-GCM',
            iv: encryptedData.iv
        },
        key,
        encryptedData.data
    );
    return dec.decode(decryptedContent);
}

async function savePrivateKey(privateKeyArmored, privateKeyPassword, masterPassword) {
    if (!privateKeyArmored || !masterPassword) {
        throw Error('Please provide both a private key and a master password.');
    }

    const salt = window.crypto.getRandomValues(new Uint8Array(16));
    const key = await deriveKey(masterPassword, salt);
    const encryptedPrivateKey = await encryptData(key, privateKeyArmored, privateKeyPassword);

    localStorage.setItem('privateKey', JSON.stringify({
        salt: Array.from(salt),
        iv: Array.from(encryptedPrivateKey.iv),
        data: Array.from(encryptedPrivateKey.data)
    }));
}

function clearPrivateKey() {
    localStorage.removeItem('privateKey');
    document.getElementById('privateKey').value = '';
    alert('Private key cleared.');
    document.getElementById('showPasswordBtn').disabled = true;
}

async function decryptPassword(encryptedContent, privateKeyArmored, privateKeyPassword) {

    try {
        const privateKey = await openpgp.readPrivateKey({ armoredKey: privateKeyArmored });

        const decryptedKey = await openpgp.decryptKey({
            privateKey,
            passphrase: privateKeyPassword
        });

        const encryptedMessage = await openpgp.readMessage({
            binaryMessage: encryptedContent
        });

        const { data: decrypted } = await openpgp.decrypt({
            message: encryptedMessage,
            decryptionKeys: decryptedKey,
        });

        return decrypted;
    } catch (error) {
        console.error('Decryption failed:', error);
        throw error;
    }
}

async function handlePasswordDecrypt(privateKey, privateKeyPassword, masterPassword, encryptedContent) {
    if (!masterPassword) {
        throw Error('Please enter your master password.');
    }

    const binaryString = atob(encryptedContent);
    const len = binaryString.length;
    const uint8Array = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        uint8Array[i] = binaryString.charCodeAt(i);
    }

    var decryptedPrivateKeyArmored = privateKey;
    if (!privateKey){
        try{
            const storedPrivateKeyData = JSON.parse(localStorage.getItem('privateKey'));
            if (!storedPrivateKeyData) {
                throw Error('No private key found. Please save your private key first.');
            }

            const salt = new Uint8Array(storedPrivateKeyData.salt);
            const key = await deriveKey(masterPassword, salt);
            let decryptedData = await decryptData(key, {
                iv: new Uint8Array(storedPrivateKeyData.iv),
                data: new Uint8Array(storedPrivateKeyData.data)
            });
            const obj = JSON.parse(decryptedData);
            decryptedPrivateKeyArmored = obj.privateKey;
            privateKeyPassword = obj.privateKeyPassword;
            console.log(decryptedPrivateKeyArmored);
        } catch (error) {
            throw error;
        }
    }

    const decryptedContent = await decryptPassword(uint8Array, decryptedPrivateKeyArmored, privateKeyPassword);
    return decryptedContent;
}
