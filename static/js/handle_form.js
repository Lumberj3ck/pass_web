
async function signAndSubmit(privateKeyArmored){
    evt.preventDefault(); // prevent default htmx request
    const messageText = 'Hello, this is a signed message.';
    let authInput = document.getElementById("auth-private-key")


    let passphrase = "test";
    const privateKey = await openpgp.readPrivateKey({ armoredKey: privateKeyArmored });
    const decryptedKey = await openpgp.decryptKey({
        privateKey,
        passphrase
    });


    console.log(decryptedKey)
    const signedMessage = await openpgp.sign({
        message: await openpgp.createMessage({ text: messageText }),
        signingKeys: decryptedKey
    });

    document.getElementById("signature").value = signedMessage
    console.log(signedMessage)
    htmx.trigger("#signed-form", "submit");
}
