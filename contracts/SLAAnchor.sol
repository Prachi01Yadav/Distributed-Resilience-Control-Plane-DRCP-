// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract SLAAnchor {
    event BreachRecorded(string serviceId, bytes32 merkleRoot, uint256 timestamp);

    struct Breach {
        string serviceId;
        bytes32 merkleRoot;
        uint256 timestamp;
    }

    Breach[] public breaches;
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can record breaches");
        _;
    }

    function recordBreach(string memory _serviceId, bytes32 _merkleRoot) public onlyOwner {
        Breach memory newBreach = Breach({
            serviceId: _serviceId,
            merkleRoot: _merkleRoot,
            timestamp: block.timestamp
        });

        breaches.push(newBreach);
        emit BreachRecorded(_serviceId, _merkleRoot, block.timestamp);
    }

    function getBreachesCount() public view returns (uint256) {
        return breaches.length;
    }
}
